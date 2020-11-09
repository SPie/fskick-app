<?php

namespace App\Console;

use Symfony\Component\Console\Helper\Table;
use Symfony\Component\Console\Output\OutputInterface;

/**
 * Class TableHelperFactory
 *
 * @package App\Console
 */
class TableHelperFactory
{
    /**
     * @param OutputInterface $output
     * @param array           $headers
     * @param array           $rows
     *
     * @return Table
     */
    public function create(OutputInterface $output, array $headers, array $rows): Table
    {
        return (new Table($output))
            ->setHeaders($headers)
            ->setRows($rows);
    }
}
